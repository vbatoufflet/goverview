/*jshint
    browser: true,
    devel: true,
    jquery: true,
    trailing: true,
    unused: true
*/

/*globals
    angular
*/

var app = angular.module('goverview', [
    'cfp.hotkeys',
    'ngCookies',
    'ngResource'
]);

app.config(['hotkeysProvider', function (hotkeysProvider) {
    hotkeysProvider.includeCheatSheet = false;
}]);

app.factory('api', ['$resource', function ($resource) {
    return $resource('/api/:type/:name?', {
        type: '@type',
        name: '@name'
    }, {
        get: {
            method: 'GET',
            params: {
                type: '@type',
                name: '@name'
            }
        },
        list: {
            method: 'GET',
            params: {
                type: '@type',
                name: null
            },
            isArray: true
        },
        search: {
            method: 'POST',
            params: {
                type: 'search',
                name: null
            },
            isArray: true
        }
    });
}]);

app.controller("MainController", ['$scope', '$cookies', '$interval', '$timeout', 'hotkeys', 'api', function ($scope,
    $cookies, $interval, $timeout, hotkeys, api) {

    $scope.loading = true;

    $scope.data = [];
    $scope.nodes = [];
    $scope.groups = [];
    $scope.services = [];

    $scope.filter = '';
    $scope.filterShow = false;
    $scope.tooltip = {};
    $scope.tooltipShow = false;
    $scope.paneShow = false;

    function getHigherState(services) {
        var state = -1;

        angular.forEach(services, function (service) {
            if (service.state > state)
                state = service.state;
        });

        return state;
    }

    $scope.getService = function (host, name) {
        var result = null;

        angular.forEach(host.services, function (service) {
            if (service.name !== name)
                return;

            result = service;
        });

        return result;
    };

    $scope.refresh = function () {
        api.list({
            type: 'nodes'
        }, function (data) {
            $scope.nodes = data;
        });

        api.list({
            type: 'groups'
        }, function (data) {
            $scope.groups = data;
        });

        api.search($scope.query, function (data) {
            var dispatch = {},
                keys,
                i,
                hosts = [],
                services = [],
                timeRef,
                count = 0;

            // Get “new” state time delta
            timeRef = moment().subtract(1, 'minute');

            // Parse retrieved data
            angular.forEach(data, function (host) {
                var state;

                host.state_new = moment(host.state_changed).isAfter(timeRef);

                // Dispatch hosts by state criticity
                state = host.state !== 0 ? host.state + 3 : getHigherState(host.services);

                if (!dispatch[state])
                    dispatch[state] = [];

                dispatch[state].push(host);

                // Get distinct services names
                angular.forEach(host.services, function (service) {
                    service.state_new = moment(service.state_changed).isAfter(timeRef);

                    if (services.indexOf(service.name) == -1)
                        services.push(service.name);
                });

                // Update states count
                count += host.state !== 0 ? 1 : host.services.length;
            });

            $scope.services = services;

            // Sort data by criticity
            keys = Object.keys(dispatch);
            keys.sort().reverse();

            for (i in keys)
                hosts = hosts.concat(dispatch[keys[i]]);

            $scope.data = hosts;
            $scope.services.sort();

            // Update document title
            document.title = document.title.replace(/(\s+\([0-9]+\))?$/, ' (' + count + ')');

            $scope.error = false;
            $scope.loading = false;
        }, function () {
            $scope.services = [];
            $scope.data = [];

            $scope.error = true;
            $scope.loading = false;
        });
    };

    $scope.resetForm = function () {
        $scope.form = {
            nodes: [],
            groups: [],
            states: [],
            acknowledges: false,
            downtimes: false
        };

        $scope.query = null;
    };

    $scope.setTooltip = function (e, host, service) {
        if (e.type == 'mouseleave') {
            if (e.relatedTarget === $scope.tooltip.target)
                return;
            else if ($scope.tooltipTimeout)
                $timeout.cancel($scope.tooltipTimeout);

            if ($(e.relatedTarget).closest('.tooltip').length === 0) {
                $scope.tooltipShow = false;
                $scope.tooltip.target = null;
            }

            return;
        }

        if (e.target === $scope.tooltip.target)
            return;
        else if ($scope.tooltipTimeout)
            $timeout.cancel($scope.tooltipTimeout);

        $scope.tooltipTimeout = $timeout(function () {
            var $target,
                obj,
                time,
                offset;

            // Update tooltip information
            obj = service ? service : host;
            time = moment(obj.state_changed);

            $scope.tooltipShow = true;
            $scope.tooltip.target = e.target;
            $scope.tooltip.changed = time.format('lll');
            $scope.tooltip.changed_relative = time.fromNow();
            $scope.tooltip.node = host.node;
            $scope.tooltip.acknowledges = obj.acknowledges;
            $scope.tooltip.comments = obj.comments;
            $scope.tooltip.links = obj.links;
            $scope.tooltip.output = service ? service.output : '';

            // Set tooltip position
            $target = $(e.target);

            offset = $target.offset();

            $('#tooltip').css({
                top: offset.top,
                left: $target.closest('td.host').length > 0 ? e.clientX : offset.left,
            });
        }, 500);
    };

    $scope.$watch('filter', function (newValue, oldValue) {
        if (newValue === oldValue)
            return;

        if ($scope.filterTimeout)
            $timeout.cancel($scope.filterTimeout);

        $scope.filterTimeout = $timeout(function () {
            $scope.form.filter = newValue;
        }, 500);
    });

    $scope.$watch('form', function () {
        $scope.query = angular.extend({}, $scope.form);

        if ($scope.query.states.length > 0)
            $scope.query.states = $scope.query.states.map(function (a) { return parseInt(a, 10); });
        else
            $scope.query.states = [1, 2, 3];

        // Save options
        $cookies.putObject('opts', $scope.form);

        $scope.refresh();
    }, true);

    // Register hotkeys
    hotkeys.add({
        combo: 'esc',
        allowIn: ['INPUT'],
        callback: function (e) {
            e.preventDefault();

            if (!e.target.value)
                $timeout(function () { e.target.blur(); }, 0);
            else
                $scope.filter = '';
        }
    });

    hotkeys.add({
        combo: 'enter',
        allowIn: ['INPUT'],
        callback: function (e) {
            if (e.target.tagName != 'INPUT')
                return;

            $timeout(function () { e.target.blur(); }, 0);
        }
    });

    hotkeys.add({
        combo: '/',
        callback: function (e) {
            e.preventDefault();
            $scope.filterShow = true;
            $timeout(function () { $('input[name=filter]').select(); }, 0);
        }
    });

    hotkeys.add({
        combo: 'p',
        callback: function () {
            $scope.paneShow = !$scope.paneShow;
        }
    });

    // Perform form initialization
    $scope.form = $cookies.getObject('opts');
    if ($scope.form === undefined)
        $scope.resetForm();
    else if ($scope.form.filter)
        $scope.filter = $scope.form.filter;

    // Set refresh interval
    $interval($scope.refresh, 30000);
}]);
