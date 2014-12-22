app.controller('NewUserController', ['$scope', '$routeParams', '$http', '$location', function($scope, $routeParams, $http, $location) {
  $scope.user = {};

  $http.post("/user/new")
}]);
