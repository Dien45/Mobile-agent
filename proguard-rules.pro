# Proguard rules for release build
-keepattributes Signature
-keepattributes *Annotation*
-keep class kotlin.** { *; }
-keep class androidx.** { *; }
-dontwarn kotlinx.**
