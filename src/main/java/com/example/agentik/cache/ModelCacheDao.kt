package com.example.agentik.cache

import androidx.room.Dao
import androidx.room.Insert
import androidx.room.OnConflictStrategy
import androidx.room.Query
import com.example.agentik.cache.CachedModel

@Dao
interface ModelCacheDao {
    @Insert(onConflict = OnConflictStrategy.REPLACE)
    suspend fun insert(model: CachedModel)

    @Query("SELECT * FROM cached_models WHERE provider_id = :providerId ORDER BY fetched_at DESC LIMIT 1")
    suspend fun getBestModelForProvider(providerId: String): CachedModel?

    @Query("DELETE FROM cached_models WHERE fetched_at < :threshold")
    suspend fun cleanupOldModels(threshold: Long)
}
