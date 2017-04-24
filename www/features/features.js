'use strict';

angular.module('myApp.features', ['ngRoute', 'angular-timeline', 'btford.markdown'])

.config(['$routeProvider', function($routeProvider) {
  $routeProvider.when('/features', {
    templateUrl: 'features/featuresTimeline.html',
    controller: 'FeaturesCtrl'
  });
}])
.config(['markdownConverterProvider', function (markdownConverterProvider) {
  // options to be passed to Showdown
  // see: https://github.com/coreyti/showdown#extensions
  markdownConverterProvider.config({
    extensions: ['github']
  });
}])

.controller('FeaturesCtrl', function(Features, $scope, $timeout, $routeParams) {
  $scope.curated = $routeParams.curate || false;
  $scope.dockerImportantFeatures = [];

  $scope.start = function() {
      Features.getFeatures({curate: $routeParams.curate}).success(function(data, status) {
          if (status == 200) {
            for (var i = 0; i < data.length; i++) {
              data[i].merged_at = moment(data[i].merged_at).fromNow();
              data[i].created_at = moment(data[i].created_at).fromNow();
              $scope.dockerImportantFeatures.push(data[i]);
            }
            $scope.show = true;
          }
      });
  }

  $scope.start();


  $scope.curate = function(pr, status) {
      Features.curate(pr, status);
      $scope.dockerImportantFeatures.forEach(function(f, i) {
          if (f.number == pr.number) {
              $scope.dockerImportantFeatures.splice(i, 1);
          }
      });
  }

  //infinite-scroll
  $scope.getNextPage = function() {
    if ($scope.dockerImportantFeatures.length == 0) {
        return;
    }
    var searchParams = {
        query: "",
        curate: $routeParams.curate,
        skip: $scope.dockerImportantFeatures[$scope.dockerImportantFeatures.length - 1].merged_at,
    }
    if ($scope.loadingFeatures == true) {
      return;
    }
    $scope.loadingFeatures = true;

    Features.getFeatures(searchParams).success(function (data, status) {
          if (status == 200) {
            for (var i = 0; i < data.length; i++) {
              data[i].merged_at = moment(data[i].merged_at).fromNow();
              data[i].created_at = moment(data[i].created_at).fromNow();
              $scope.dockerImportantFeatures.push(data[i]);
            }
            // Make sure the layout is rendered before enabling again
            $timeout(function() {$scope.loadingFeatures = false}, 500);
            return;
          }
    });
  };
});
