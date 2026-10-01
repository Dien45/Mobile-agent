package com.example.agentik.di

import android.content.Context
import androidx.room.Room
import com.example.agentik.ai.HermesAgent
import com.example.agentik.cache.ModelCacheDao
import com.example.agentik.provider.Provider
import com.example.agentik.provider.PresetProvider
import com.example.agentik.util.NetworkUtil
import kotlinx.coroutines.CoroutineScope
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.SupervisorJob

object ServiceLocator {
    // Application context set during init
    private lateinit var appContext: Context
    private val appScope = CoroutineScope(SupervisorJob() + Dispatchers.IO)

    // Simple mutable provider (default preset OPENAI)
    var currentProvider: Provider = Provider(
        id = "openai",
        name = "OpenAI",
        baseUrl = "https://api.openai.com/v1",
        apiKey = "",
        isPreset = true,
        weight = 1
    )
        private set

    // Initialize with application context
    fun init(context: Context) {
        appContext = context.applicationContext
        // Initialize Room database
        database = Room.databaseBuilder(
            appContext,
            AppDatabase::class.java,
            "agentik-db"
        ).fallbackToDestructiveMigration().build()
        // Initialize network util (holds base URL via provider)
        networkUtil = NetworkUtil()
        // Initialize HermesAgent
        hermesAgent = HermesAgent(currentProvider, database.modelCacheDao())
    }

    // Room database instance
    private lateinit var database: AppDatabase
    fun getModelCacheDao(): ModelCacheDao = database.modelCacheDao()

    // Network utility
    lateinit var networkUtil: NetworkUtil
        private set

    // Hermes agent instance (recreated when provider changes)
    lateinit var hermesAgent: HermesAgent
        private set

    // Update provider (called from SettingsScreen)
    fun updateProvider(provider: Provider) {
        currentProvider = provider
        // Recreate dependent components
        hermesAgent = HermesAgent(currentProvider, database.modelCacheDao())
    }
}
