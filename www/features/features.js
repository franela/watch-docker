'use strict';

angular.module('myApp.features', ['ngRoute', 'angular-timeline', 'btford.markdown'])

.config(['$routeProvider', function($routeProvider) {
  $routeProvider.when('/features', {
    templateUrl: 'features/featuresTimeline.html',
    controller: 'FeaturesCtrl',
    resolve: {
          features: function(Features) {
            return Features.getFeatures();
          }
    }
  });
}])
.config(['markdownConverterProvider', function (markdownConverterProvider) {
  // options to be passed to Showdown
  // see: https://github.com/coreyti/showdown#extensions
  markdownConverterProvider.config({
    extensions: ['github']
  });
}])

.controller('FeaturesCtrl', function(features, Features, $scope, $timeout) {

  $scope.dockerImportantFeatures = features.data;
  $scope.show = true;

  //infinite-scroll
  $scope.getNextPage = function() {
    var searchParams = {
        query: "",
        skip: $scope.dockerImportantFeatures[$scope.dockerImportantFeatures.length - 1].merged_at,
    }
    if ($scope.loadingFeatures == true) {
      return;
    }
    $scope.loadingFeatures = true;

    Features.getFeatures(searchParams).success(function (data, status) {
          if (status == 200) {
            for (var i = 0; i < data.length; i++) {
              $scope.dockerImportantFeatures.push(data[i]);
            }
            // Make sure the layout is rendered before enabling again
            $timeout(function() {$scope.loadingFeatures = false}, 500);
            return;
          }
    });
  };
});
