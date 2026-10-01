package com.example.agentik.di

import android.content.Context
import com.example.agentik.ai.HermesAgent
import com.example.agentik.cache.ModelCacheDao
import com.example.agentik.provider.Provider
import com.example.agentik.util.NetworkUtil
import org.koin.android.ext.koin.androidContext
import org.koin.core.context.startKoin
import org.koin.dsl.module

fun initKoin(context: Context) {
    startKoin {
        androidContext(context)
        modules(appModule)
    }
    // Initialize ServiceLocator (which also creates DB, network, hermes)
    ServiceLocator.init(context)
}

private val appModule = module {
    // Provide Provider (mutable via ServiceLocator)
    single { ServiceLocator.currentProvider }
    // Provide ModelCacheDao from AppDatabase
    single { ServiceLocator.getModelCacheDao() }
    // Provide NetworkUtil (holds base URL etc.)
    single { ServiceLocator.networkUtil }
    // Provide HermesAgent (recreated on provider change via ServiceLocator)
    single { ServiceLocator.hermesAgent }
}
