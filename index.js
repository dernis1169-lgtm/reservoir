const { Server } = require("socket.io");

const io = new Server({
  maxHttpBufferSize: 1e7, // Увеличиваем лимит для передачи фото (10MB)
  cors: {
    origin: "*",
    methods: ["GET", "POST"]
  }
});

io.on("connection", (socket) => {
  console.log("User connected:", socket.id);

  socket.on("joinRoom", (roomId) => {
    socket.join(roomId);
    console.log(`User ${socket.id} joined room: ${roomId}`);
  });

  socket.on("draw", (data) => {
    if (data && data.roomId) {
      socket.to(data.roomId).emit("draw", data);
    }
  });

  // Событие очистки экрана
  socket.on("clear", (data) => {
    if (data && data.roomId) {
      socket.to(data.roomId).emit("clear");
    }
  });

  // Событие смены фонового фото
  socket.on("setBackground", (data) => {
    if (data && data.roomId) {
      socket.to(data.roomId).emit("setBackground", data.image);
    }
  });

  socket.on("reaction", (data) => {
    if (data && data.roomId) {
      socket.to(data.roomId).emit("reaction", data);
    }
  });

  socket.on("disconnect", () => {
    console.log("User disconnected:", socket.id);
  });
});

const PORT = process.env.PORT || 8080;
io.listen(PORT);
console.log(`Socket.IO Server running on port ${PORT}`);
