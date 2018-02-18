# python script for attacker probabilities

import math


def nCr(n,r):
    a = 1
    for i in range(r):
        a *= n - i
    return a / math.factorial(r)


def attack_success_probability(attackers_in_pool, attackers_in_review, pool_size, reviewers):
    return nCr(attackers_in_pool, attackers_in_review) * nCr(pool_size - attackers_in_pool, reviewers - attackers_in_review) / nCr(pool_size, reviewers)

def g(attackers_in_pool, attackers_in_review):
    return attack_success_probability(attackers_in_pool, attackers_in_review, 100, 5)

def get_data():
    for i in range(5, 95):
        attacker_probabilities = [g(i,2),g(i, 3), g(i, 4), g(i, 5)]
    #    attacker_probabilities = [g(i, 1), g(i, 2), g(i, 3), g(i, 4), g(i, 5)]
        attacker_probabilities_string = [str(p) for p in attacker_probabilities]
        print('%d,%s,%s' % (i, ','.join(attacker_probabilities_string), sum(attacker_probabilities)))
