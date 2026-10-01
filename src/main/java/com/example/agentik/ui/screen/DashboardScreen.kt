package com.example.agentik.ui.screen

import androidx.compose.foundation.layout.*
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.verticalScroll
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import androidx.navigation.NavController
import com.example.agentik.ai.HermesAgent
import com.example.agentik.navigation.Screen
import kotlinx.coroutines.flow.collectLatest
import kotlinx.coroutines.launch
import org.koin.androidx.compose.get

@Composable
fun DashboardScreen(navController: NavController) {
    var showShell by remember { mutableStateOf(false) }
    val shellOutput = remember { mutableStateListOf<String>() }
    var userInput by remember { mutableStateOf("") }
    val hermesAgent: HermesAgent = get()
    val coroutineScope = rememberCoroutineScope()

    Column(
        modifier = Modifier
            .fillMaxSize()
            .padding(24.dp),
        verticalArrangement = Arrangement.Top
    ) {
        Text("Dashboard", style = MaterialTheme.typography.headlineMedium)
        Spacer(modifier = Modifier.height(16.dp))
        if (showShell) {
            // Output area
            Column(
                modifier = Modifier
                    .weight(1f)
                    .verticalScroll(rememberScrollState())
            ) {
                shellOutput.forEach { line ->
                    Text(line)
                }
            }
            Spacer(modifier = Modifier.height(8.dp))
            Row {
                OutlinedTextField(
                    value = userInput,
                    onValueChange = { userInput = it },
                    label = { Text("Perintah") },
                    modifier = Modifier.weight(1f)
                )
                Spacer(modifier = Modifier.width(8.dp))
                Button(onClick = {
                    if (userInput.isNotBlank()) {
                        // Show user command
                        shellOutput.add("▶ $userInput")
                        // Invoke HermesAgent chat and stream response
                        coroutineScope.launch {
                            hermesAgent.chat(userInput).collectLatest { chunk ->
                                shellOutput.add(chunk)
                            }
                        }
                        userInput = ""
                    }
                }) {
                    Text("Kirim")
                }
            }
        } else {
            Button(onClick = { showShell = true }) {
                Text("Buka Shell (TODO)")
            }
        }
        Spacer(modifier = Modifier.height(8.dp))
        Button(onClick = { navController.navigate(Screen.Login.route) }) {
            Text("Logout")
        }
    }
}
