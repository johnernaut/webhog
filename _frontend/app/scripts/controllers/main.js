'use strict';

angular.module('webhog')
  .controller('MainCtrl', [
    '$scope',
    'entities',
    'Restangular',
    function($scope, entities, Restangular) {
      $scope.entities = entities;

      $scope.delete = function(e) {
        Restangular.one('entity', e.id).remove().then(function() {
          $scope.entities.splice($scope.entities.indexOf(e), 1);
        }, function(data) {
          console.log('error');
        });
      };
    }
  ]);
