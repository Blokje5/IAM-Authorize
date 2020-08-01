package main

test_has_action_matches_action {
    has_action[["policyid", _]]
    
    with data.policies.policyid.statements as {
		{"actions": ["iam:CreateUser"]},
		{"actions": ["iam:CreateRole"]},
	}

    with input.action as "iam:CreateUser"
}

test_has_resource_matches_resource {
    has_resource[["policyid", _]]
    
    with data.policies.policyid.statements as {
		{"resources": ["iam:User1"]},
		{"resources": ["iam:User2"]},
	}

    with input.resource as "iam:User2"
}

test_allow {
    allow

    with data.policies as {
        {
            "statements": [{
                "effect": "allow",
                "actions": ["iam:CreateUser"],
                "resources": ["iam:User1"]
            }]
        }
    }

    with input as {
        "action": "iam:CreateUser",
        "resource": "iam:User1"
    }
}

test_deny {
    deny

    with data.policies as {
        {
            "statements": [{
                "effect": "deny",
                "actions": ["iam:CreateUser"],
                "resources": ["iam:User1"]
            }]
        }
    }

    with input as {
        "action": "iam:CreateUser",
        "resource": "iam:User1"
    }
}

test_authorized_simple_allow {
    authorized

    with data.policies as {
        {
            "statements": [{
                "effect": "allow",
                "actions": ["iam:CreateUser"],
                "resources": ["iam:User1"]
            }]
        }
    }

    with input as {
        "action": "iam:CreateUser",
        "resource": "iam:User1"
    }
}

test_authorized_simple_deny {
    not authorized

    with data.policies as {
        {
            "statements": [{
                "effect": "deny",
                "actions": ["iam:CreateUser"],
                "resources": ["iam:User1"]
            }]
        }
    }

    with input as {
        "action": "iam:CreateUser",
        "resource": "iam:User1"
    }
}
