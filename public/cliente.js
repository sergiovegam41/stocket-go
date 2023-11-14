document.addEventListener('DOMContentLoaded', function() {
    // Asegúrate de que conectas al servidor con el protocolo y puerto correcto.
    // Si estás sirviendo tu página desde el mismo servidor y puerto, puedes usar window.location.origin
    var socket = io(window.location.origin);

    socket.on('connect', function() {
        console.log('Conectado al servidor!');
        // Emite cualquier evento inicial que necesites aquí
        // socket.emit('tuEvento', { tusDatos });
    });

    // Escucha por los mensajes del servidor y actualiza el DOM
    socket.on('server:refresh:mensajes', function(data) {
        console.log(data); // Asegúrate de que los datos se reciben correctamente
        var messages = document.getElementById('messages');
        var li = document.createElement('li');
        li.textContent = typeof data === 'string' ? data : JSON.stringify(data);
        messages.appendChild(li);
    });

    window.submitMessage = function() {
        var messageInput = document.getElementById('messageInput');
        var message = messageInput.value;
        socket.emit('client:send:message', { mensaje: message });
        messageInput.value = ''; // Limpia el input después de enviar
    };
});
