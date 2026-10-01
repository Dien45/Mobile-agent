package com.example.agentik.network

import com.example.agentik.provider.Provider
import okhttp3.RequestBody
import okhttp3.ResponseBody
import retrofit2.Call
import retrofit2.http.Body
import retrofit2.http.GET
import retrofit2.http.Header
import retrofit2.http.POST
import retrofit2.http.Url

interface ApiService {
    // Dynamic base URL per request (full URL passed in)
    @GET
    fun getModels(@Url url: String, @Header("Authorization") auth: String): Call<ResponseBody>

    @POST
    fun postChat(@Url url: String, @Header("Authorization") auth: String, @Body body: RequestBody): Call<ResponseBody>
}
