// var tooltipTriggerList = [].slice.call(
//     document.querySelectorAll('[data-toggle="tooltip"]')
// );
//
// var tooltipList = tooltipTriggerList.map(function (tooltipTriggerEl) {
//     return new bootstrap.Tooltip(tooltipTriggerEl);
// });

var alertList = document.querySelectorAll('.alert')
alertList.forEach(function (alert) {
    new bootstrap.Alert(alert)
})

function enableTooltips() {
    var tooltipTriggerList = [].slice.call(document.querySelectorAll('[data-bs-toggle="tooltip"]'))
    var tooltipList = tooltipTriggerList.map(function (tooltipTriggerEl) {
      return new bootstrap.Tooltip(tooltipTriggerEl, {trigger : 'hover'})
    })
}
enableTooltips();

function disableTooltips() {
    var myTooltipElements = document.querySelectorAll('[data-bs-toggle="tooltip"]');

    for (let i=0; i < myTooltipElements.length; i++) {
        myTooltipElements[i].hide();//.tooltip("destroy");
        //tooltip.hide();
    }
    //myTooltipElements.addEventListener('hidden.bs.tooltip', function () {
      // do something...
    //})
}

//Initialize bootstrap tooltips
// var tooltipTriggerList = [].slice.call(document.querySelectorAll('[data-bs-toggle="tooltip"]'));
// var tooltipList = tooltipTriggerList.map( function(tooltipTriggerEl) {
//   return new bootstrap.Tooltip(tooltipTriggerEl, {
//   trigger : 'hover'
//   });
// });
