package com.example.agentik.ai

import com.example.agentik.cache.ModelCacheDao
import com.example.agentik.provider.Provider
import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.callbackFlow

class HermesAgent(
    private val provider: Provider,
    private val modelCache: ModelCacheDao
) {
    private val SYSTEM_PROMPT = "Kamu adalah Hermes Agent, asisten CLI Alpine Linux."

    suspend fun chat(prompt: String): Flow<String> = callbackFlow {
        val model = modelCache.getBestModelForProvider(provider.id) ?: "gpt-3.5-turbo"

        val payload = mapOf(
            "model" to model,
            "messages" to listOf(
                mapOf("role" to "system", "content" to SYSTEM_PROMPT),
                mapOf("role" to "user", "content" to prompt)
            ),
            "stream" to true
        )

        // TODO: Implementasi HTTP request ke provider.baseUrl + "/v1/chat/completions"
        // Parse SSE stream → kirim ke collector via trySend()
        trySend("Hermes response untuk: $prompt")
        close()
    }
}