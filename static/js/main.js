function connect(path) {
    var ws = new WebSocket("ws://localhost:3000/ws/"+path);
    var list = $("#todo-list");

    ws.onmessage = function(e) {
        var data = JSON.parse(e.data);

        if (data.old_val === null && data.new_val !== null) {
            // new item
            var item = data.new_val;
            list.append(""+
                "<li data-id='"+item.id+"' class='"+item.Status+"'>"+
                    "<div class='view'>"+
                        "<a href='/toggle/"+item.id+"' class='button toggle'></a>"+
                        "<span>"+item.Text+"</span>"+
                        "<a href='/delete/"+item.id+"' class='button destroy'></a>"+
                    "</div>"+
                "</li>"+
            "");
        } else if (data.old_val !== null && data.new_val === null) {
            // deleted item
            var item = data.old_val;
            $("[data-id='"+item.id+"']").remove();
        } else {
            // updated item
            var item = data.new_val;
            $("[data-id='"+item.id+"']").attr("class", item.Status);
        }
    };
}
