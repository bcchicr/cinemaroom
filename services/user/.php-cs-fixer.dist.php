<?php

$finder = (new PhpCsFixer\Finder())
    ->in(__DIR__)
    ->exclude('var')
    ->exclude('vendor')
;

return (new PhpCsFixer\Config())
    ->setRules([
        '@PER-CS3.0' => true,
        '@Symfony' => true,

        'array_syntax' => ['syntax' => 'short'],
        'trailing_comma_in_multiline' => ['elements' => ['arrays', 'arguments', 'parameters']],
        'no_unused_imports' => true,
        'ordered_imports' => ['sort_algorithm' => 'alpha'],
        'no_useless_else' => true,
        'no_useless_return' => true,
        'method_chaining_indentation' => true,
        'blank_line_before_statement' => [
            'statements' => ['break', 'continue', 'declare', 'return', 'throw', 'try'],
        ],
        'phpdoc_indent' => true,
        'single_line_comment_style' => ['comment_types' => ['hash']],
        'visibility_required' => ['elements' => ['property', 'method', 'const']],
    ])
    ->setFinder($finder)
;
