app.controller('NewDeckController', ['$scope', '$routeParams', '$http', '$location', function($scope, $routeParams, $http, $location){
  console.log($scope);

  $http.defaults.headers.common['Authorization'] = "Bearer " + localStorage.token;  

  $scope.admin = true;

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

    setTimeout(function() {
      cradQuantities = document.querySelectorAll(".crad-quantity");
      cradQuantities[cradQuantities.length - 1].focus();
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
  }

  $scope.handleCradChange = function(crad, e) {
    console.log("KEY HIT BABY!")
    if (e.keyCode === 13) {
      // handle enter keypress
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