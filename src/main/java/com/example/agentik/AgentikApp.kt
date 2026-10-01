package com.example.agentik

import android.app.Application
import com.example.agentik.di.initKoin

class AgentikApp : Application() {
    override fun onCreate() {
        super.onCreate()
        // Initialize Koin DI container
        initKoin(this)
    }
}
