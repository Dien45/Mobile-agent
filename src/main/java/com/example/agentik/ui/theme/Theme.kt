package com.example.agentik.ui.theme

import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.darkColorScheme
import androidx.compose.material3.lightColorScheme
import androidx.compose.runtime.Composable
import androidx.compose.ui.graphics.Color

val PrimaryBlue = Color(0xFFB3C7F7)
val SurfaceLight = Color(0xFFF5F7FF)

val DarkColorScheme = darkColorScheme(
    primary = PrimaryBlue,
    surface = SurfaceLight,
    onPrimary = Color(0xFF1A1A2E),
    onSurface = Color(0xFF1A1A2E)
)

val LightColorScheme = lightColorScheme(
    primary = PrimaryBlue,
    surface = SurfaceLight,
    onPrimary = Color(0xFF1A1A2E),
    onSurface = Color(0xFF1A1A2E)
)

@Composable
fun AgentikTheme(darkTheme: Boolean = false, content: @Composable () -> Unit) {
    val colors = if (darkTheme) DarkColorScheme else LightColorScheme
    MaterialTheme(colorScheme = colors, typography = Typography, content = content)
}
</arg_value>