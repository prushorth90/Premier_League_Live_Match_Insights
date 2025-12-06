package com.example.myapplication

package com.example.pitchside

import androidx.compose.foundation.background
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.unit.dp
import okhttp3.*
import okio.ByteString

// 1. WebSocket Client Setup (OkHttp)
class MatchWebSocketListener(val onMessageReceived: (String) -> Unit) : WebSocketListener() {
    override fun onMessage(webSocket: WebSocket, text: String) {
        onMessageReceived(text)
    }
}

@Composable
fun MatchLiveScreen() {
    // State to hold messages
    val matchEvents = remember { mutableStateListOf<String>() }
    val client = OkHttpClient()

    // 2. Connect to Go Backend on Launch
    LaunchedEffect(Unit) {
        val request = Request.Builder().url("ws://10.0.2.2:8080/live-updates").build() // 10.0.2.2 is localhost for Android Emulator
        val listener = MatchWebSocketListener { msg ->
            matchEvents.add(msg)
        }
        client.newWebSocket(request, listener)
    }

    Column(
        modifier = Modifier
            .fillMaxSize()
            .padding(16.dp)
    ) {
        Text(
            text = "🔴 Live Match Feed",
            style = MaterialTheme.typography.headlineMedium,
            modifier = Modifier.padding(bottom = 16.dp)
        )

        // 3. Display Events in a List
        LazyColumn(
            modifier = Modifier
                .fillMaxSize()
                .background(Color(0xFFF5F5F5))
                .padding(8.dp)
        ) {
            items(matchEvents) { event ->
                EventCard(event)
            }
        }
    }
}

@Composable
fun EventCard(eventText: String) {
    Card(
        modifier = Modifier
            .fillMaxWidth()
            .padding(vertical = 4.dp),
        colors = CardDefaults.cardColors(containerColor = Color.White),
        elevation = CardDefaults.cardElevation(defaultElevation = 2.dp)
    ) {
        Text(
            text = eventText,
            modifier = Modifier.padding(16.dp),
            style = MaterialTheme.typography.bodyLarge
        )
    }
}