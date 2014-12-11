var app = angular.module("crad", ['ngRoute', 'ngTouch']);

app.config([
  "$routeProvider", "$locationProvider", function($routeProvider, $locationProvider) {
    $routeProvider.when("/app/login", {
      templateUrl: "/assets/html/login.html",
      controller: "LoginController"
    }).when("/app/account/new", {
      templateUrl: "/assets/html/new-user.html",
      controller: "NewUserController"      
    }).when("/app/account/:username", {
      templateUrl: "/assets/html/account.html",
      controller: "AccountController"      
    }).when("/app/decks/:username/new", {
      templateUrl: "/assets/html/new-deck.html",
      controller:  "NewDeckController"      
    }).when("/app/decks/:username/:deckname", {
      templateUrl: "/assets/html/deck.html",
      controller:  "DeckController"
    }).when("/app/decks/:username/:deckname/edit", {
      templateUrl: "/assets/html/edit-deck.html",
      controller:  "EditDeckController"      
    }).otherwise({
      templateUrl: "/assets/html/homepage.html",
      controller:  "HomeController"      
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

app.directive("navMask", function(){
  return {
    restrict: "E",
    template: "<div class='mask' ng-if='active' ng-click='toggleNav()'></div>",
    controller: "NavMaskController"
  }
});

app.directive("navAside", function() {
  return {
    restrict: "E",
    templateUrl: "/assets/html/nav-aside.html",
    controller: "NavAsideController"
  }
});

app.controller('HomeController', ['$rootScope', '$scope', '$location', function($rootScope, $scope, $location) {
  $rootScope.homepage = true;
  $rootScope.$on("$routeChangeSuccess", function() {
    if ($location.path() === "/" || $location.path() === "" || $location.path() === "/app/") {
      $rootScope.homepage = true;
    } else {
      $rootScope.homepage = false;

    }
  })
}]);

app.controller('NavMaskController', ['$rootScope', '$scope', function($rootScope, $scope) {
  $scope.toggleNav = function() {
    $rootScope.active = !$rootScope.active;
  }
}]);

app.controller('SitewideHeaderController', ['$scope', '$location', function($scope, $location){
  $scope.active   = false;
  $scope.dropdown = false;
  $scope.toggleHeader = function() {
    $scope.active = !$scope.active;
  }

  $scope.$on("$routeChangeSuccess", function() {
    $scope.dropdown = '';
  });

  $scope.isActive = function(path) {
    // console.log($location.path().indexOf(path) > -1);
    return $location.path().indexOf(path) > -1;
  }

  $scope.toggleDropdown = function(what) {
    if ($scope.dropdown === what) { $scope.dropdown = ''; return }
    $scope.dropdown = what;
  }
  $scope.activeDropdown = function(what) {
    return $scope.dropdown === what;
  }

}]);

app.controller('NavAsideController', ['$scope', '$rootScope', '$http', '$location',function($scope, $rootScope, $http, $location){
  $scope.active = false;
  if (localStorage.token) {
    $http.defaults.headers.common['Authorization'] = "Bearer " + localStorage.token;
    $rootScope.loggedIn = true;
  } else {
    $rootScope.loggedIn = false;
  }

  if (localStorage.user) {
    $scope.user = JSON.parse(localStorage.user);
    $http.post("/admin/" + $scope.user.username, {}).success(function(data) {
      $scope.user = data;
    }).error(function() {
      $rootScope.loggedIn = false;
    });
  } else {
    $scope.user = {};
  }

  $scope.$on("$routeChangeSuccess", function() {
    $rootScope.active = false;
  });

  $rootScope.$on("toggleheader", function() {
    $scope.active = !$scope.active;
  });

  $rootScope.$on("deckadded", function(e, deckWrapper) {
    // console.log(arguments);
    $scope.user.decks.push(deckWrapper.deck);
  });

  // console.log($scope.user.decks);

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
    // console.log("function is running dawg")
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

  $scope.hasCrads = function() {
    // console.log($scope.deck)
    for (var key in $scope.deck.crads) {
      if (hasOwnProperty.call($scope.deck.crads, key)) return true;
    }
  }
  // if ($rootScope.user.)
  $http.get("/decks/" + $routeParams.username + "/" + $routeParams.deckname)
    .success(function(data) {
      // console.log("SUCCESSSS");
      // console.log(data);
      $scope.deck = data;
      for (cradName in $scope.deck.crads) {
        var crad = $scope.deck.crads[cradName];
        // console.log(crad);
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
            // console.log(data);
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

app.controller('AccountController', ['$scope', '$http', '$routeParams', function($scope, $http, $routeParams){
  $http.post("/admin/" + $routeParams.username, {}).success(function(data) {
    $scope.user = data;
  });  
}]);




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


app.controller('EditDeckController', ['$scope', '$routeParams', '$http', '$location', function($scope, $routeParams, $http, $location){
  console.log($scope);

  $http.defaults.headers.common['Authorization'] = "Bearer " + localStorage.token;  

  $scope.admin = false;

  $scope.deck = {
    crads: []
  };

  if ($routeParams.deckname) {
    $http.post("/admin/" + $routeParams.username, {}).success(function(data) {
      $scope.admin = true;
      $http.get("/decks/" + $routeParams.username + "/" + $routeParams.deckname)
        .success(function(data) {
          console.log(data);
          $scope.deck = data;

          oldCrads = $scope.deck.crads;
          $scope.deck.crads = [];
          for (crad in oldCrads) {
            $scope.deck.crads.push(oldCrads[crad]);
          }
          $scope.deck.crads.push({quantity:1, name: ""});
          // console.log($scope.deck.crads);
        });
    });
  }

  $scope.suggestCard = function(crad) {
    if (crad.name && crad.name.length > 0) {
      $http.get("/crads/" + crad.name)
        .success(function(data) {
          crad.suggestions = data;
          console.log(data);
        });

    }
  }

  $scope.setName = function(crad, suggestion, deck) {
    deck.crads[crad.name] = {};
    crad.name = suggestion;
    deck.crads[crad.name] = crad;
    delete crad.suggestions;
    deck.crads["new-crad"] = {quantity:null, name:""};

    // get the last one

    // setTimeout(function() {
    //   // cradQuantities = document.querySelectorAll(".crad-quantity");
    //   // cradQuantities[cradQuantities.length - 1].focus();
    // }, 10);
    if ($scope.deck.crads[$scope.deck.crads.length-1] === crad)  {
      // if this crad is the last crad! 
      $scope.deck.crads.push({quantity: 1, name: ""});
    }
    setTimeout(function() {
      setTimeout(function() {
        for (var crad in $scope.deck.crads) {
          delete $scope.deck.crads[crad].suggestions;
        }
        $scope.$apply();
      }, 10)
      var qs = document.querySelectorAll(".quantity")
      for (var i = 0; i < qs.length; qs++) {
        qs[i].value.trim()
      }

      e.target.parentElement.nextElementSibling.children[0].focus();
    }, 10);

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
    if (e.keyCode === 13) {
      return false;
    }
  }

  $scope.preventSubmit = function(e) {
    if (e.keyCode === 13) {
      e.preventDefault();
      return false;      
    }
  }

  $scope.handleCradChange = function(crad, e, i) {
    console.log("KEY HIT BABY!");
    e.preventDefault();
    if (e.keyCode === 13) {
      // handle enter keypress
      if ($scope.deck.crads[$scope.deck.crads.length-1] === crad)  {
        // if this crad is the last crad! 
        $scope.deck.crads.push({quantity: 1, name: ""});
      }
      setTimeout(function() {
        setTimeout(function() {
          for (var crad in $scope.deck.crads) {
            delete $scope.deck.crads[crad].suggestions;
          }
          $scope.$apply();
        }, 10)
        var qs = document.querySelectorAll(".quantity")
        for (var i = 0; i < qs.length; qs++) {
          qs[i].value.trim()
        }

        e.target.parentElement.nextElementSibling.children[0].focus();

      }, 10);
    } else if (e.keyCode === 8) {
      // handle backspace
      console.log(crad.name);
      if (crad.name) {
        if(crad.name.length === 0 || e.target.selectionEnd === 0 ) {
          setTimeout(function() {
            e.target.previousElementSibling.focus();
          }, 0);
        }
      }
    }
    return false;
  }

  $scope.deckHasError = function(errorList) {
    // console.log("Errors! ", errors === {});
    for (errors in errorList) {
      for (var i = 0; i < errorList[errors].length; i++) {
        console.log()
        if (errorList[errors][i].$touched) {
          return true
        }
      }
    }
    return false;
  }

  $scope.saveDeck = function() {
    console.log("SAVING!");
    delete $scope.deck.crads["new-crad"];
    cradArray = [];
    for(var crad = 0; crad < $scope.deck.crads.length; crad++) {
      cradToAdd = $scope.deck.crads[crad]
      if (cradToAdd.name !== undefined && cradToAdd.name.length > 1 ) {
        cradArray.push(cradToAdd);
      }
    }
    console.log("Crads:", $scope.deck.crads);
    console.log("Crad Array:", cradArray);
    $http.put("/decks/" + $routeParams.username, {
      deck:  {
        "name":  $routeParams.deckname,
        "crads": cradArray
      }
      
    }).success(function(data) {
      // woo hoo it works!
      console.log("Your deck has been saved dawg");
      $location.path("/app/decks/" + $routeParams.username + "/" + $routeParams.deckname);
    });
  }

}]);

app.controller('NewUserController', ['$scope', '$routeParams', '$http', '$location', function($scope, $routeParams, $http, $location) {
  
}]);


