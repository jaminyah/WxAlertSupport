
var captchaForm = {
        CaptchaType: "character",
        Id: '',
        VerifyValue: '',
        ConfigAudio: {
            CaptchaLen: 6,
            Language: 'en'
        },
        ConfigCharacter: {
            Height: 60,
            Width: 240,
            Mode: 2,
            ComplexOfNoiseText: 0,
            ComplexOfNoiseDot: 0,
            IsUseSimpleFont: true,
            IsShowHollowLine: false,
            IsShowNoiseDot: false,
            IsShowNoiseText: false,
            IsShowSlimeLine: false,
            IsShowSineLine: false,
            CaptchaLen: 6,
            BgColor:{R:3, G:102, B:214, A:255}
        },
        ConfigDigit: {
            Height: 80,
            Width: 240,
            CaptchaLen: 5,
            MaxSkew: 0.7,
            DotCount: 80
        }
    }

$('#comment-form').submit(function(e) {

    e.preventDefault();         // avoid executing the actual submit form
    document.getElementById("loader").style.display = "block";

    var form = $(this);
    $.ajax({
        type: form.attr('method'),
        url: "/submit",
        contentType: 'application/x-www-form-urlencoded',
        data: form.serialize(),
        success: function(data) {
            document.getElementById("loader").style.display = "none";
            $("#username").val("");
            $("#message").val("");
            $("#captcha-solution").val("");
            var comments = JSON.parse(data);
            displayComments(comments);
            generateCaptcha();
            $(".captcha-row").show(2000);
        },
        error: function(data) {
            console.log('There is an error');
            console.log(data);
        }
    });
});

$(document).ready(function() {
   $("#captcha-solution").val("");
   $("#message").val("");

    $.ajax({
        type: "get",
        url: "/comments",
        success: function(data) {
            document.getElementById("loader").style.display = "none";
            var comments = JSON.parse(data);
            displayComments(comments);
            generateCaptcha();
        },
        error: function(data) {
            console.log('There is an error');
            console.log(data);
        }
    });
});


function displayComments(json_data) {
    console.log(json_data);
    var comment = "";
    var list = $("<ul>");

    for (var i = 0; i < json_data.length; i++) {
        comment = "<img src='./images/avatar.png'/>" +
        "<h3>" + json_data[i]['name'] + "</h3>" +
        "<h4>" + json_data[i]['date'] + "</h4>" +
        "<p>" + json_data[i]['comment'] + "</p>";
            
        var item = $("<li>").html(comment);
        list.append(item);
    }
    $("#comment-area").html(list);
}
    
function generateCaptcha() {
    console.log('Generating captcha');
    const url = 'api/getCaptcha';
    var blob = "";

    let fetchData = {
        method: 'post',
        body: JSON.stringify(captchaForm),
        headers: new Headers()
    }

    fetch(url, fetchData)
    .then(function(response){
        console.log(" fetch .then");
        return response.json();
    })
    .then(function(data){
        console.log(data.captchaId);
        captchaForm.Id = data.captchaId;
        blob = data.data;
        displayCaptcha(data);
    })
    .catch(function(error){
        console.log("fetch error: ")
        console.log(error);
    });
}

function displayCaptcha(captcha) {
    let captchaImage = "<p><h3>Verify you are human</h3>" + "<h4>Solve the math problem:</h4></p>" + "<img src='" + captcha.data + "'/>";
    $("#captcha-img").html(captchaImage);
}

function verifyCaptcha() {
    console.log("verifyCaptcha");
    captchaForm.VerifyValue = document.getElementById("captcha-solution").value;

    const url = 'api/verifyCaptcha';

    let fetchData = {
        method: 'post',
        body: JSON.stringify(captchaForm),
        headers: new Headers()
    }

    fetch(url, fetchData)
    .then(function(response){
        return response.json();
    })
    .then(function(data){
        console.log(data.msg);
        if (data.code == "success") {
            $(".captcha-row").hide("slow", function(){
                showMessage(data.msg);
            });
        } else {
            console.log(data.code);
            showMessage(data.msg);
            generateCaptcha();
        }
    })
    .catch(function(error){
        console.log("fetch error: ")
        console.log(error);
    });
}

function showMessage(msgText) {
    console.log("Show message.")
    let msg = msgText;
    $(".alert").find('.message').text(msg);
    $(".alert").fadeIn("slow", function() {
        setTimeout(function(){
            $(".alert").fadeOut("slow");
        }, 2000);
    });
}
