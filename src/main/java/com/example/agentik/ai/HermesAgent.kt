package com.example.agentik.ai

import com.example.agentik.provider.Provider
import com.example.agentik.cache.ModelCacheDao
import com.example.agentik.network.ApiService
import com.example.agentik.util.Logger
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.flow
import kotlinx.coroutines.withContext
import okhttp3.MediaType.Companion.toMediaType
import okhttp3.RequestBody.Companion.toRequestBody
import org.json.JSONObject

class HermesAgent(
    private val provider: Provider,
    private val modelCacheDao: ModelCacheDao,
    private val apiService: ApiService? = null // injected via Koin later
) {
    /**
     * Get best model for the current provider from cache, fallback to a default.
     */
    private suspend fun getModelId(): String = withContext(Dispatchers.IO) {
        modelCacheDao.getBestModelForProvider(provider.id)?.modelId ?: "gpt-3.5-turbo"
    }

    /**
     * Fetch model list from provider and cache it.
     */
    suspend fun refreshModels() = withContext(Dispatchers.IO) {
        try {
            val url = "${provider.baseUrl}/v1/models"
            val authHeader = "Bearer ${provider.apiKey}"
            val call = apiService?.getModels(url, authHeader)
            val response = call?.execute()
            if (response != null && response.isSuccessful) {
                val bodyStr = response.body?.string() ?: return@withContext
                val json = JSONObject(bodyStr)
                // Support both OpenAI style ("data") and custom ("models")
                val modelsArray = if (json.has("data")) json.getJSONArray("data") else json.getJSONArray("models")
                for (i in 0 until modelsArray.length()) {
                    val obj = modelsArray.getJSONObject(i)
                    val modelId = obj.getString("id")
                    val modelName = obj.optString("name", modelId)
                    val ownedBy = obj.optString("owned_by", "unknown")
                    modelCacheDao.insert(com.example.agentik.cache.CachedModel(
                        providerId = provider.id,
                        modelId = modelId,
                        modelName = modelName,
                        ownedBy = ownedBy,
                        fetchedAt = System.currentTimeMillis()
                    ))
                }
                Logger.i("HermesAgent", "Models refreshed for ${provider.name}")
            } else {
                Logger.w("HermesAgent", "Failed to fetch models: ${response?.code()}")
            }
        } catch (e: Exception) {
            Logger.e("HermesAgent", "Exception while refreshing models", e)
        }
    }

    /**
     * Send a chat prompt to the provider and receive a **Flow** of streamed responses.
     * The flow emits raw text chunks (SSE "data:" lines) as they arrive.
     */
    fun chat(prompt: String): Flow<String> = flow {
        val modelId = runCatching { getModelId() }.getOrElse { "gpt-3.5-turbo" }
        val url = "${provider.baseUrl}/v1/chat/completions"
        val authHeader = "Bearer ${provider.apiKey}"
        // Build OpenAI‑compatible payload (system + user messages)
        val payload = JSONObject().apply {
            put("model", modelId)
            put("stream", true)
            put("messages", listOf(
                mapOf("role" to "system", "content" to "You are Hermes, an AI assistant.") ,
                mapOf("role" to "user", "content" to prompt)
            ))
        }
        val requestBody = payload.toString().toRequestBody("application/json".toMediaType())
        try {
            val call = apiService?.postChat(url, authHeader, requestBody)
            val response = call?.execute()
            if (response != null && response.isSuccessful) {
                val source = response.body?.source() ?: return@flow
                // Simple SSE parser – each line starting with "data:" is a chunk
                while (!source.exhausted()) {
                    val line = source.readUtf8Line() ?: break
                    if (line.startsWith("data:")) {
                        val data = line.removePrefix("data:").trim()
                        if (data == "[DONE]") break
                        // Expect OpenAI style {"choices":[{"delta":{"content":"text"}}]}
                        try {
                            val json = JSONObject(data)
                            val choices = json.getJSONArray("choices")
                            if (choices.length() > 0) {
                                val delta = choices.getJSONObject(0).optJSONObject("delta")
                                val content = delta?.optString("content")
                                if (!content.isNullOrEmpty()) {
                                    emit(content)
                                }
                            }
                        } catch (e: Exception) {
                            // If parsing fails, just emit raw line
                            emit(data)
                        }
                    }
                }
                Logger.i("HermesAgent", "Chat stream completed")
            } else {
                Logger.w("HermesAgent", "Chat request failed: ${response?.code()}")
                emit("[Error] Unable to get response from provider.")
            }
        } catch (e: Exception) {
            Logger.e("HermesAgent", "Exception during chat", e)
            emit("[Exception] ${e.message}")
        }
    }
}
