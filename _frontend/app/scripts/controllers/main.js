'use strict';

angular.module('webhog')
  .controller('MainCtrl', [
    '$scope',
    'entities',
    function ($scope, entities) {
      $scope.entities = entities;
    }
  ]);
