package postgres

// @note: each model will also use Init, as constructors are not idiomatic go
//        and may create a problem in that because we need to deal with errors
//        any constructor would not make parent code any cleaner.
//
//        since each model being managed is different we cannot use a single
//        interface as an abstraction, but we can standardize the operations
//        so that each interface will be mostly the same.
