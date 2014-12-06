app.controller('NewDeckController', ['$scope', '$routeParams', '$http', '$location', '$rootScope', function($scope, $routeParams, $http, $location, $rootScope) {
  $http.defaults.headers.common['Authorization'] = "Bearer " + localStorage.token;  

  $scope.deck = {};
  $scope.admin = false;
  $http.post("/admin/" + $routeParams.username, {})
    .success(function(data) {
      $scope.admin = true;
    });
  $scope.saveDeck = function() {
    console.log("$scope.deck = ", $scope.deck);
    $http.post("/decks/"+ $routeParams.username, {
      deck: $scope.deck
    }).success(function(data) {
      // window.theData = data;
      // $rootScope.user.decks.push(data);
      $rootScope.$broadcast("deckadded", {
        "deck": data
      });
      $location.path("/decks/" + $routeParams.username + "/" + $scope.deck.name);
    });
  }

}]);
