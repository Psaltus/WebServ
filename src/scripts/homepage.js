$(document).ready(function(){
    $("button").on("click", function(){
        $("button").after("<b>Hello!</b> ");
     });

     $("h1").on("click", function() {
         $("#mainbody").slideToggle();
     });
});