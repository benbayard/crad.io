var app = angular.module("crad", ['ngRoute', 'ngTouch', 'ngResource', 'ngAnimate']);

app.config([
  "$routeProvider", "$locationProvider", function($routeProvider, $locationProvider) {
    $routeProvider.when("/app/login", {
      templateUrl: "/assets/html/login.html",
      controller: "LoginController"
    }).when("/app/account", {
      templateUrl: "/assets/html/account.html",
      controller: "AccountController"      
    }).when("/app/decks/:username/:deckname", {
      templateUrl: "/assets/html/deck.html",
      controller:  "DeckController"
    }).when("/app/decks/:username/new", {
      templateUrl: "/assets/html/new-deck.html",
      controller:  "NewDeckController"      
    }).when("/app/decks/:username/:deckname/edit", {
      templateUrl: "/assets/html/new-deck.html",
      controller:  "NewDeckController"      
    });
    
    return $locationProvider.html5Mode(true);
  }
]);

app.directive("sitewideHeader", function() {
  return {
    restrict: "E",
    templateUrl: "/assets/html/sitewide-header.html",
    controller: "SitewideHeaderController"
  };
});

app.directive("navAside", function() {
  return {
    restrict: "E",
    templateUrl: "/assets/html/nav-aside.html",
    controller: "NavAsideController"
  }
});

app.controller('SitewideHeaderController', ['$scope', function($scope){
  
}]);

app.controller('NavAsideController', ['$scope', '$rootScope', '$http', '$location',function($scope, $rootScope, $http, $location){
  if (localStorage.token) {
    $http.defaults.headers.common['Authorization'] = "Bearer " + localStorage.token;
    $rootScope.loggedIn = true;
  } else {
    $rootScope.loggedIn = false;
  }

  if (localStorage.user) {
    $scope.user = JSON.parse(localStorage.user);
  } else {
    $scope.user = {};
  }

  console.log($scope.user.decks);

  $scope.isActive = function(path) {
    return path == $location.path();
  }

  $scope.logOut = function() {
    delete localStorage.token
    $rootScope.loggedIn = false;
  }
}]);

app.controller('LoginController', ['$scope', '$http', '$rootScope', function($scope, $http, $rootScope){
  $scope.user = {};
  $scope.getToken = function() {
    console.log("function is running dawg")
    $http.post("/login", {
      "email":    $scope.user.email,
      "password": $scope.user.password
    }).success(function(data) {
      $rootScope.user   = data.user;
      localStorage.user = JSON.stringify($rootScope.user)
      $rootScope.token  = localStorage.token = data.token;
      $rootScope.loggedIn = true;
    }).error(function(data) {

    });
  }
}]);

app.controller('DeckController', ['$scope', '$http', '$routeParams', '$rootScope', '$location', function($scope, $http, $routeParams, $rootScope, $location){
  $http.defaults.headers.common['Authorization'] = "Bearer " + localStorage.token;

  $scope.deck = {};
  $scope.active = 'column';
  $scope.types = {};

  $scope.editing = false;

  $scope.startEditing = function() {
    // $scope.editing = true;
    $location.path("/app/decks/" + $routeParams.username + "/" + $routeParams.deckname + "/edit");
  }

  $scope.isActive = function(type) {
    return type == $scope.active;
  }
  // if ($rootScope.user.)
  $http.get("/decks/" + $routeParams.username + "/" + $routeParams.deckname)
    .success(function(data) {
      console.log("SUCCESSSS");
      console.log(data);
      $scope.deck = data;
      for (cradName in $scope.deck.crads) {
        var crad = $scope.deck.crads[cradName];
        console.log(crad);
        window.crads = $scope.deck.crads;
        if(localStorage[crad.name]) {
          $scope.deck.crads[crad.name].cradData = JSON.parse(localStorage[crad.name]);
          if ($scope.types[crad.cradData.types[0]] == undefined) {
            $scope.types[crad.cradData.types[0]] = new Array($scope.deck.crads[crad.name]);
          } else {
            $scope.types[crad.cradData.types[0]].push($scope.deck.crads[crad.name])
          }
        } else {
          $http.get("/crad/" + crad.name).success(function(data) {
            localStorage[data.name] = JSON.stringify(data);
            $scope.deck.crads[data.name].cradData = data;
            console.log(data);
            window.crads = $scope.deck.crads;
            if ($scope.types[crad.cradData.types[0]] == undefined) {
              $scope.types[crad.cradData.types[0]] = new Array($scope.deck.crads[crad.name]);
            } else {
              $scope.types[crad.cradData.types[0]].push($scope.deck.crads[crad.name])
            }
          });          
        }
      }
        
    });
  
  $http.post("/admin/" + $routeParams.username, {}).success(function(data) {
    $scope.admin = true;
  });

}]);

app.controller('AccountController', ['$scope', '$http', function($scope, $http){
  
}]);



app.controller('NewDeckController', ['$scope', '$routeParams', '$http', function($scope, $routeParams, $http){
  console.log($scope);

  $scope.admin = true;

  if ($routeParams.deckname) {
    $http.post("/admin/" + $routeParams.username, {}).success(function(data) {
      $scope.admin = true;
      $http.get("/decks/" + $routeParams.username + "/" + $routeParams.deckname)
        .success(function(data) {
          console.log("GOT THE FUCKING DECK ASSHOLE");
          $scope.deck = data;
          if($scope.deck.crads !== undefined) {
            $scope.deck.crads = [];
          }
          $scope.deck.crads.push({quantity: 0, name: ""});

          console.log($scope.deck.crads);
        });
    });
  }

  $scope.suggestCard = function(crad) {
    if (crad.name.length > 0) {
      $http.get("/crads/" + crad.name)
        .success(function(data) {
          crad.suggestions = data;
          console.log(data);
        });

    }
  }

  $scope.setName = function(crad, suggestion, deck) {
    crad.name = suggestion;
    delete crad.suggestions
    deck.crads.push({quantity:null, name:""});
    cradQuantities = document.querySelectorAll(".crad-quantity");

    // get the last one
    cradQuantities[cradQuantities.length - 1].focus();

    return true;
  }

  $scope.isCradName = function(crad, suggestion) {
    return crad.name === suggestion;
  }

  $scope.gotoCradname = function(crad, e) {
    if (e.keyCode === 32) {
      // console.log(e.target.value);
      e.target.value = e.target.value.replace(/\s/g, "");
      setTimeout(function() {
        e.target.nextElementSibling.focus();
      }, 0);
      return false;
    }
  }

  $scope.handleCradChange = function(crad, e) {
    console.log("KEY HIT BABY!")
    if (e.keyCode === 13) {
      // handle enter keypress
    } else if (e.keyCode === 8) {
      // handle backspace
      console.log(crad.name);
      if(!crad.name || crad.name.length === 0 || e.target.selectionEnd === 0 ) {
        setTimeout(function() {
          e.target.previousElementSibling.focus();
        }, 0);
      }
    }
  }

}]);

