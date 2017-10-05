$.ajaxSetup({
    beforeSend: function (xhr) {
        var token = $("meta[name='_csrf']").attr("content");
        var header = $("meta[name='_csrf_header']").attr("content");
        xhr.setRequestHeader(header, token);
    }
});

$(function () {
    $("p.markdown").each(function (e) {
        var txt = $(this).text();
        $(this).html(marked(txt));
    });

    $("form[data-confirm]").submit(function (e) {
        if (!confirm($(this).data('confirm'))) {
            e.preventDefault();
        }
    });

    $("a[data-method]").click(function (e) {
        e.preventDefault();
        var msg = $(this).data('confirm');
        var method = $(this).data('method');
        var next = $(this).data('next');
        var url = $(this).attr('href');

        var ok = true;
        if (msg) {
            if (!confirm(msg)) {
                ok = false;
            }
        }
        if (ok) {
            // console.log(method, url, next);
            $.ajax({type: method, url: url}).done(function (rst) {
                // FIXME
                if (rst.ok) {
                    window.location.href = next;
                } else {
                    alert(rst);
                }
            }).fail(function (jqXHR, textStatus, errorThrown) {
                alert(textStatus);
            });
        }
    });
});
