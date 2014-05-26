'use strict';

var app = angular
  .module('webhog', [
    'ngCookies',
    'ngResource',
    'restangular',
    'chieffancypants.loadingBar',
    'mgcrea.ngStrap',
    'ngSanitize',
    'ngRoute'
  ])
  .config(function ($routeProvider, $locationProvider) {
    $locationProvider.hashPrefix('!');
    $routeProvider
      .when('/', {
        templateUrl: 'views/main.html',
        controller: 'MainCtrl',
        resolve: {
          entities: ['$route', 'Restangular', function ($route, Restangular) {
            return Restangular.all('entities').getList();
          }]
        }
      })
      .when('/entity/:slug', {
        templateUrl: 'views/entity.html',
        controller: 'EntityController',
        resolve: {
          entity: ['$route', 'Restangular', function ($route, Restangular) {
            return Restangular.one('entity', $route.current.params.slug).get();
          }]
        }
      })
      .otherwise({
        redirectTo: '/'
      });
  });

app.run([
  'Restangular',
  '$alert',
  '$location',
  '$rootScope',
  function (Restangular, $alert, $location, $rootScope) {
    Restangular.setBaseUrl('/api');
    Restangular.setDefaultHeaders({'X-API-KEY': 'SCRAPEAPI', 'Content-Type': 'application/json'});
    return Restangular.setErrorInterceptor(function (res) {
      if (res.status === 404 || res.status === 401) {
        $location.path('/');
        $alert({
          title: res.statusText,
          content: res.data.error,
          placement: 'top-right',
          type: 'warning',
          duration: 5,
          show: true
        });
      }
      return false;
    });
  }
]);