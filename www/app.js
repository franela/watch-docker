'use strict';

// Declare app level module which depends on views, and components
angular.module('myApp', [
  'ngRoute',
  'myApp.features',
  'ui.bootstrap',
  'myApp.Factories',
  'infinite-scroll'
]).
config(['$locationProvider', '$routeProvider', function($locationProvider, $routeProvider) {
  $locationProvider.hashPrefix('!');
  $routeProvider.otherwise({redirectTo: '/features'});

}]);

var appFactories = angular.module('myApp.Factories', []);
