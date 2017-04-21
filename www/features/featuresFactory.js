appFactories
  .factory('Features', function ($http) {

  return {
          getFeatures: function(searchParams) {
            return $http.get("/timeline", {params: searchParams});
        }
    };
  });
