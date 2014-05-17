'use strict';

angular.module('webhog')
  .controller('EntityController', [
    '$scope',
    'entity',
    function ($scope, entity) {
      $scope.entity = entity;
    }
  ]);
