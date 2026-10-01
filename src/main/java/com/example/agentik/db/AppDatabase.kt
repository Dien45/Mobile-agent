package com.example.agentik.db

import androidx.room.Database
import androidx.room.RoomDatabase
import com.example.agentik.cache.CachedModel
import com.example.agentik.cache.ModelCacheDao

@Database(entities = [CachedModel::class], version = 1, exportSchema = false)
abstract class AppDatabase : RoomDatabase() {
    abstract fun modelCacheDao(): ModelCacheDao
}
