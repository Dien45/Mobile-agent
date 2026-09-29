package com.example.agentik.provider

data class Provider(
    val id: String,
    val name: String,
    val baseUrl: String,
    val apiKey: String,
    val isPreset: Boolean = true,
    val weight: Int = 1
)