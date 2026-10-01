package com.example.agentik.cache

import androidx.room.Database
import androidx.room.RoomDatabase
import com.example.agentik.cache.ModelCacheDao
import com.example.agentik.cache.CachedModel

@Database(entities = [CachedModel::class], version = 1, exportSchema = false)
abstract class AppDatabase : RoomDatabase() {
    abstract fun modelCacheDao(): ModelCacheDao
}
