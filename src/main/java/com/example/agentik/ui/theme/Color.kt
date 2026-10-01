package com.example.agentik.ui.theme

import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Typography
import androidx.compose.ui.graphics.Color

val PrimaryBlue = Color(0xFFB3C7F7)
val SurfaceLight = Color(0xFFF5F7FF)
val PastelGray = Color(0xFFAAB7B8)

val ColorScheme = androidx.compose.material3.darkColorScheme(
    primary = PrimaryBlue,
    surface = SurfaceLight,
    onPrimary = Color(0xFF1A1A2E),
    onSurface = Color(0xFF1A1A2E)
)

val LightColorScheme = androidx.compose.material3.lightColorScheme(
    primary = PrimaryBlue,
    surface = SurfaceLight,
    onPrimary = Color(0xFF1A1A2E),
    onSurface = Color(0xFF1A1A2E)
)
