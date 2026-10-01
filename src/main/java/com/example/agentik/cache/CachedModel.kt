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
    @ColumnInfo(name = "fetched_at") val fetchedAt: Long = System.currentTimeMillis()
)
