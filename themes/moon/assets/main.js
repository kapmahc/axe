$.ajaxSetup({
  beforeSend: function(xhr) {
    xhr.setRequestHeader('Authenticity-Token', $('meta[name=X-CSRF-Token]').attr('content'));
  }
});

$(function() {
  $("p.markdown").each(function(e) {
    var txt = $(this).text();
    $(this).html(marked(txt));
  });

  $("form[data-next]").submit(function(e) {
    e.preventDefault();
    var msg = $(this).data('confirm');
    var method = $(this).attr('method');
    var next = $(this).data('next');
    var action = $(this).attr('action');
    var data = $(this).serialize();

    var ok = true;
    if (msg) {
      if (!confirm(msg)) {
        ok = false;
      }
    }
    if (ok) {
      // console.log(method, next, action, data);
      $.ajax({type: method, data: data, url: action}).done(function(rst) {
        console.log(rst)
        if (rst.message) {
          alert(rst.message);
        }
        if (next) {
          window.location.href = next;
        }
      }).fail(function(xhr) {
        alert(xhr.responseText);
      });
    }
  });

  $("a[data-method]").click(function(e) {
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
      $.ajax({type: method, url: url}).done(function(rst) {
        if (rst.message) {
          alert(rst.message);
        }
        if (next) {
          window.location.href = next;
        }
      }).fail(function(xhr) {
        alert(xhr.responseText);
      });
    }
  });
});
