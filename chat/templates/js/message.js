$(function(){
    var socket = null;
    var msgBox = $("#chatbox textarea");
    var messages = $("#messages");
    $("#chatbox").submit(function(){
        if (!msgBox.val()) return false;
        if (!socket) {
            alert("Error: There is no socket connection.");
            return false;
        }
        //SOCKET SEND
        socket.send(JSON.stringify({"Message": msgBox.val()}));

        msgBox.val("");
        return false;
    });
    if (!window["WebSocket"]) {
        alert("Error: Your browser does not support web sockets.")
    } else {
        socket = new WebSocket("ws://localhost:8080/room");
        socket.onclose = function() {
            alert("Connection has been closed.");
        }

        //Message combination append
        socket.onmessage = function(e) {
            var msg = JSON.parse(e.data)
            messages.append(
                $("<img>")
                    .attr("src", msg.AvatarURL)
                    .addClass("avatar1")
                    .attr("title", msg.Name)
                    .css({
                        width: 50,
                        height: 50
                    }),
                    $("<span>").text(msg.Name),
                    $("<div>").text(msg.Message)
             );
        }
    }
});