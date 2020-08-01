package main

import data.policies

default authorized = false

has_resource[[policy_id, statement_id]] {
    resource := policies[policy_id].statements[statement_id].resources[_]
    resource == input.resource
}

action_matches(in, stored) {
	in == stored
}

has_action[[policy_id, statement_id]] {
    statement_action := policies[policy_id].statements[statement_id].actions[_]
    action_matches(input.action, statement_action)
}

match[[effect, policy_id, statement_id]] {
    effect := policies[policy_id].statements[statement_id].effect
    has_resource[[policy_id, statement_id]]
	has_action[[policy_id, statement_id]]
}

allow {
    match[["allow", _, _]]
}

deny {
    match[["deny", _, _]]
}

authorized {
    allow
    not deny
}