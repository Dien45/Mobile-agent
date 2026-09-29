package com.example.agentik.cache

import androidx.room.Dao
import androidx.room.Insert
import androidx.room.OnConflictStrategy
import androidx.room.Query

@Dao
interface ModelCacheDao {
    @Insert(onConflict = OnConflictStrategy.REPLACE)
    suspend fun insert(model: CachedModel)

    @Query("SELECT model_id FROM cached_models WHERE provider_id = :providerId ORDER BY fetched_at DESC LIMIT 1")
    suspend fun getBestModelForProvider(providerId: String): String?

    @Query("DELETE FROM cached_models WHERE provider_id = :providerId AND fetched_at < :threshold")
    suspend fun cleanupOldModels(providerId: String, threshold: Long)
}

data class CachedModel(
    val providerId: String,
    val modelId: String,
    val modelName: String,
    val ownedBy: String,
    val fetchedAt: Long
)