package com.example.agentik.util

import android.util.Log

object Logger {
    private const val TAG = "Agentik"
    fun d(tag: String, msg: String) = Log.d("$TAG-$tag", msg)
    fun i(tag: String, msg: String) = Log.i("$TAG-$tag", msg)
    fun w(tag: String, msg: String) = Log.w("$TAG-$tag", msg)
    fun e(tag: String, msg: String, tr: Throwable? = null) {
        Log.e("$TAG-$tag", msg, tr)
    }
}
