package com.example.agentik.util

import com.example.agentik.provider.Provider
import com.example.agentik.network.ApiService
import okhttp3.OkHttpClient
import okhttp3.logging.HttpLoggingInterceptor
import retrofit2.Retrofit
import retrofit2.converter.moshi.MoshiConverterFactory

object NetworkUtil {
    fun createRetrofit(provider: Provider): ApiService {
        val logging = HttpLoggingInterceptor().apply {
            level = HttpLoggingInterceptor.Level.BODY
        }
        val client = OkHttpClient.Builder()
            .addInterceptor(logging)
            .addInterceptor { chain ->
                val original = chain.request()
                val requestBuilder = original.newBuilder()
                    .header("Authorization", "Bearer ${'$'}{provider.apiKey}")
                    .header("Content-Type", "application/json")
                chain.proceed(requestBuilder.build())
            }
            .build()
        return Retrofit.Builder()
            .baseUrl(provider.baseUrl) // baseUrl is required but will be overridden by @Url in ApiService
            .client(client)
            .addConverterFactory(MoshiConverterFactory.create())
            .build()
            .create(ApiService::class.java)
    }
}
