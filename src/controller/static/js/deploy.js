'use strict';

(function (document, window, $) {

// @include jQuery

function verifyActionCode (code) {
  // TODO
  return true;
}

function _nop() {}

function _api(path, method, data, headers) {
  var config = {
    url: '/api/v1' + path,
    method: method,
    dataType: 'json'
  };
  if (headers) {
    config.headers = headers;
  }
  if (data) {
    config.contentType = 'application/json';
    config.data = JSON.stringify(data);
  }
  return $.ajax(config);
}

var api = {
  get: function (path, data, headers) {
    return _api(path, 'GET', data, headers);
  },
  post: function (path, data, headers) {
    return _api(path, 'POST', data, headers);
  },
  put: function (path, data, headers) {
    return _api(path, 'PUT', data, headers);
  },
  delete: function (path, data, headers) {
    return _api(path, 'DELETE', data, headers);
  },
  patch: function (path, data, headers) {
    return _api(path, 'PATCH', data, headers);
  },
  custom: function (path, method, data, headers) {
    return _api(path, method, data, headers);
  }
};

function klassGroup () {
  this.raw = {
    id: -1,
    name: null,
    description: null,
    owner_id: -1,
    owner_obj: null,
    status: null,
    deployment: null,
    process: null
  };
}

function klassNode () {
  this.raw = {
    id: -1,
    name: null,
    uuid: null,
    description: null,
    status: null,
    owner_id: -1,
    owner_obj: null,
    tags: null,
    nics: null
  };
}

function klassRepo () {
  this.raw = {
    id: -1,
    owner_id: -1,
    owner_obj: null,
    name: null,
    is_public: false,
    description: null,
    readme: null,
    tags: [] // {id, name, yml}
  }
}

function klassUser () {
  this.raw = {
    id: -1,
    name: null,
    displayName: null,
    key: null, // store in cookie
    email: null,
    createTime: null
  }
}

klassUser.prototype = {
  signup: function (username, password, email, code) {
    if (!verifyActionCode()) {
      // TODO alert
      return;
    }
    return api.post('/signup', {
      username: username,
      password: password,
      email: email
    });
  },
  login: function (username, password, code) {
    if (!verifyAcitionCode()) {
      // TODO alert
      return;
    }
    return api.post('/login', {
      username: username,
      password: password,
    });
  },
  profile: function () {
    // TODO check current user
    return api.post('/user/', {
      password: password,
    });
  },
  resetPassword: function (password) {
    // TODO check current user
    return api.post('/user/reset-password', {
      password: password,
    });
  },
  resetKey: function () {
    // TODO check current user
    return api.post('/user/reset-key');
  },

  repoctl: {
    list: function () {
      return api.get('/user/repos');
    },
    create: function (is_public, name, description, readme) {
      return api.post('/user/repos', {
        name: name,
        description: description,
        readme: readme,
        is_public: is_public
      });
    },
    one: function (reponame) {
      return api.get('/user/repos/' + reponame);
    }
  },
  nodectl: {
    list: function () {
      return api.get('/user/nodes');
    },
    one: function (nodename) {
      return api.get('/user/nodes/' + nodename);
    }
  },
  groupctl: {
    list: function () {
      return api.get('/user/groups');
    },
    create: function (name, description) {
      return api.post('/user/groups', {
        name: name,
        description: description
      });
    },
    one: function (groupname) {
      return api.get('/user/groups/' + groupname);
    }
  }
};

klassRepo.prototype = {
  update: function (is_public, description, readme) {
    return api.put('/user/repos/' + this.raw.name, {
      is_public: is_public,
      description: description,
      readme: readme
    });
  },
  delete: function () {
    return api.delete('/user/repos/' + this.raw.name);
  },
  listTags: function () {
    return api.get('/user/repos/' + this.raw.name + '/tags');
  },
  addTag: function (tagname, yml) {
    return api.post('/user/repos/' + this.raw.name + '/tags', {
      name: tagname,
      yml: yml
    });
  },
  removeTag: function (tagname) {
    return api.delete('/user/repos/' + this.raw.name + '/tags/' + tagname);
  }
};

klassNode.prototype = {
  update: function (name, description) {
    return api.put('/user/nodes/' + this.raw.name, {
      name: name,
      description: description
    });
  },
  delete: function () {
    return api.delete('/user/nodes/' + this.raw.name);
  },
  listTags: function () {
    return api.get('/user/nodes/' + this.raw.name + '/tags');
  },
  addTag: function (tagname) {
    return api.post('/user/nodes/' + this.raw.name + '/tags', {
      name: tagname
    });
  },
  removeTag: function (tagname) {
    return api.delete('/user/nodes/' + this.raw.name + '/tags/' + tagname);
  },
  addNicTag: function (nicname, tagname) {
    return api.post('/user/nodes/' + this.raw.name + '/nics/' + nicname +'/tags', {
      name: tagname
    });
  },
  removeTag: function (nicname, tagname) {
    return api.delete('/user/nodes/' + this.raw.name + '/nics/' + nicname + '/tags/' + tagname);
  }
};

klassGroup.prototype = {
  update: function (name, description) {
    return api.put('/user/groups/' + this.raw.name, {
      name: name,
      description: description
    });
  },
  delete: function () {
    return api.delete('/user/groups/' + this.raw.name);
  },
  listNode: function () {
    return api.get('/user/groups/' + this.raw.name + '/nodes');
  },
  addNode: function (nodename) {
    return api.post('/user/groups/' + this.raw.name + '/nodes', {
      name: nodename
    });
  },
  removeNode: function (nodename) {
    return api.delete('/user/groups/' + this.raw.name + '/nodes/' + nodename);
  },
  listDeployments: function () {
    return api.get('/user/groups/' + this.raw.name + '/deployment');
  },
  createDeployment: function (repo) {
    return api.post('/user/groups/' + this.raw.name + '/deployment', repo.raw);
  },
  deleteDeployment: function () {
    return api.delete('/user/groups/' + this.raw.name + '/deployment');
  },
  prepareDeployment: function () {
    // TODO check status
    return api.put('/user/groups/' + this.raw.name + '/deployment/prepare');
  },
  deploy: function () {
    // TODO check status
    return api.put('/user/groups/' + this.raw.name + '/deployment/execute');
  },
  progress: function () {
    return api.get('/user/groups/' + this.raw.name + '/deployment/process');
  }
};

// exports
var moonlegend = {
  Group: klassGroup,
  Node: klassNode,
  Repo: klassRepo,
  User: klassUser
};

window.MoonLegend = moonlegend;

})(document, window, jQuery);
