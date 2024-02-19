package subscription

// success testcases
// case: create subscription with addon
// case: set cancelAtPeriodEnd subscription and billing cycle effected, and check upgrade|downgrade will resume it
// case: upgrade subscription with addon
// case: billing cycle without pendingUpdate and check dunning time invoice
// case: downgrade subscription with addon
// case: billing cycle with pendingUpdate and check dunning time invoice
// case: set subscription trialEnd and billing cycle effected, check trialEnd radius, should after max(now,periodEnd) -- todo set time not may cause sub new cycle invoice and payment
// case: upgrade|downgrade subscription after periodEnd and before trialEnd
// case: cancel subscription immediately

// failure testcases
// case1: create subscription with payment failure and check expired cycle
// case2: billing cycle with payment failure after periodEnd, sub should change to incomplete, else may set trialEnd
// case3: incomplete status situations todo
