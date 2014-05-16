'use strict';

angular.module('webhog')
  .controller('MainCtrl', [
    '$scope',
    'entity',
    function ($scope, entity) {
      $scope.entity = entity;
    }
  ]);
