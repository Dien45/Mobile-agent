package com.example.agentik.cache

import androidx.room.ColumnInfo
import androidx.room.Entity
import androidx.room.PrimaryKey

@Entity(tableName = "cached_models")
data class CachedModel(
    @PrimaryKey(autoGenerate = true) val uid: Int = 0,
    @ColumnInfo(name = "provider_id") val providerId: String,
    @ColumnInfo(name = "model_id") val modelId: String,
    @ColumnInfo(name = "model_name") val modelName: String,
    @ColumnInfo(name = "owned_by") val ownedBy: String,
    @ColumnInfo(name = "fetched_at") val fetchedAt: Long
)

package com.example.agentik.cache

import androidx.room.Dao
import androidx.room.Insert
import androidx.room.OnConflictStrategy
import androidx.room.Query

@Dao
interface ModelCacheDao {
    @Insert(onConflict = OnConflictStrategy.REPLACE)
    fun insert(cachedModel: CachedModel)

    @Query("SELECT * FROM cached_models WHERE provider_id = :providerId ORDER BY fetched_at DESC LIMIT 1")
    fun getBestModelForProvider(providerId: String): CachedModel?

    @Query("DELETE FROM cached_models WHERE fetched_at < :cutoff")
    fun cleanupOldModels(cutoff: Long)
}
