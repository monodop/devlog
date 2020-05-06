grammar FilterGrammar;
entry: filter* EOF;

filter: plaintext #filter_simple
      | equals    #filter_key
      ;
equals: plaintext EQUALS plaintext;
plaintext: WORD | STRING;

STRING: '"' ( ~('"' | '\\') | '\\' ('"' | '\\') )* '"';
EQUALS: '=';
WORD: ~([ \n\t\r=])+;
WS: [ \n\t\r]+ -> skip;