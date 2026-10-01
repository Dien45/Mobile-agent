package com.example.agentik.provider

// Provider data class representing an LLM endpoint configuration.
// Removed kotlinx.serialization.Serializable to avoid needing the serialization plugin.

data class Provider(
    val id: String,
    val name: String,
    val baseUrl: String,
    val apiKey: String,
    val isPreset: Boolean = false,
    val weight: Int = 1
)

enum class PresetProvider(
    val id: String,
    val displayName: String,
    val baseUrl: String
) {
    OPENAI("openai", "OpenAI", "https://api.openai.com"),
    GROQ("groq", "Groq", "https://api.groq.com"),
    ANTHROPIC("anthropic", "Anthropic", "https://api.anthropic.com"),
    OLLAMA("ollama", "Ollama", "http://127.0.0.1:11434"),
    OMNIRoute_TAILSCALE("omniroute", "OmniRoute (Tailscale)", "http://omniroute.local:3000")
}
