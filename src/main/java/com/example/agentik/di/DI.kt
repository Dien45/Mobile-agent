package com.example.agentik.di

import android.content.Context
import androidx.room.Room
import com.example.agentik.cache.AppDatabase
import com.example.agentik.cache.ModelCacheDao
import com.example.agentik.provider.Provider
import com.example.agentik.provider.PresetProvider
import com.example.agentik.ai.HermesAgent
import com.example.agentik.network.ApiService
import com.example.agentik.util.NetworkUtil
import org.koin.android.ext.koin.androidContext
import org.koin.core.context.startKoin
import org.koin.dsl.module

fun initKoin(context: Context) {
    startKoin {
        androidContext(context)
        modules(appModule)
    }
}

private val appModule = module {
    // Room database
    single {
        Room.databaseBuilder(
            get(),
            AppDatabase::class.java,
            "agentik-db"
        ).fallbackToDestructiveMigration().build()
    }
    // DAO
    single { get<AppDatabase>().modelCacheDao() }

    // Default provider (preset OpenAI) – can be overridden via Settings UI
    single {
        Provider(
            id = PresetProvider.OPENAI.id,
            name = PresetProvider.OPENAI.displayName,
            baseUrl = PresetProvider.OPENAI.baseUrl,
            apiKey = "",
            isPreset = true,
            weight = 1
        )
    }

    // Retrofit ApiService (dynamic URL per request)
    single<ApiService> { NetworkUtil.createRetrofit(get()) }

    // HermesAgent – core AI wrapper
    single { HermesAgent(get(), get(), get()) }
}
