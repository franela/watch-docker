appFactories
  .factory('Features', function ($http) {

  return {
          getFeatures: function(searchParams) {
            return $http.get("/timeline", {params: searchParams});
          },
          curate: function(pr, status) {
            return $http.put("/pulls", {repo: pr.base.repo.full_name, number: pr.number, curated: status});
          }
  };
});
