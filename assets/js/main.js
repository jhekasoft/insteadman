function getOs() {
    // This script sets OSName variable as follows:
    // "Windows"    for all versions of Windows
    // "MacOS"      for all versions of Macintosh OS
    // "Linux"      for all versions of Linux
    // "UNIX"       for all other UNIX flavors
    // "Unknown OS" indicates failure to detect the OS
    var osName="Unknown OS";
    if (navigator.appVersion.indexOf("Win")!=-1) osName = "windows"
    else if (navigator.appVersion.indexOf("Mac")!=-1) osName = "macosx"
    else if (navigator.appVersion.indexOf("X11")!=-1) osName = "unix"
    else if (navigator.appVersion.indexOf("Linux")!=-1) osName="gnulinux"

    return osName;
}

$(document).ready(function() {
    var osName = getOs();

    var $buttonEl;
    if ("windows" == osName) {
        $buttonEl = $(".btn-download-windows");
    } else if ("macosx" == osName) {
        $buttonEl = $(".btn-download-macosx");
    } else if ("gnulinux" == osName || "unix" == osName) {
        $buttonEl = $(".btn-download-gnulinux");
    }

    $buttonEl.removeClass("btn-default");
    $buttonEl.addClass("btn-success");
});

$(".release-info-icon, .release-info-close").click(function() {
    $(".release-info").slideToggle();
})
