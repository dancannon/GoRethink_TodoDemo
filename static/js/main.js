function connect(path) {
    var ws = new WebSocket("ws://localhost:3000/ws/"+path);
    var list = $("#todo-list");

    ws.onmessage = function(e) {
        var data = JSON.parse(e.data);

        if (data.OldValue === null && data.NewValue !== null) {
            // new item
            var item = data.NewValue;
            list.append(""+
                "<li data-id='"+item.id+"' class='"+item.Status+"'>"+
                    "<div class='view'>"+
                        "<a href='/toggle/"+item.id+"' class='button toggle'></a>"+
                        "<span>"+item.Text+"</span>"+
                        "<a href='/delete/"+item.id+"' class='button destroy'></a>"+
                    "</div>"+
                "</li>"+
            "");
        } else if (data.OldValue !== null && data.NewValue === null) {
            // deleted item
            var item = data.OldValue;
            $("[data-id='"+item.id+"']").remove();
        } else {
            // updated item
            var item = data.NewValue;
            $("[data-id='"+item.id+"']").attr("class", item.Status);
        }
    };
}
